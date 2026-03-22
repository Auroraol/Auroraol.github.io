/**
 * Sync to GitHub Skill - 优化版
 * 将 markdown 文件同步到 GitHub（使用 GitHub API，避免 git 权限问题）
 */

const fs = require('fs').promises;
const path = require('path');
const matter = require('gray-matter');
const https = require('https');

class SyncToGitHub {
  constructor(config) {
    // 优先从 SKILLS_CONFIG_JSON 中读取配置
    let skillsConfig = {};
    console.error('🔍 [SyncToGitHub] 开始初始化...');
    
    if (process.env.SKILLS_CONFIG_JSON) {
      try {
        const parsedConfig = JSON.parse(process.env.SKILLS_CONFIG_JSON);
        skillsConfig = parsedConfig['sync-to-github'] || {};
        console.error('✅ [SyncToGitHub] 从 SKILLS_CONFIG_JSON 读取配置成功');
      } catch (e) {
        console.error('❌ [SyncToGitHub] 解析 SKILLS_CONFIG_JSON 失败:', e.message);
      }
    }
    
    // 参数优先级：constructor config > SKILLS_CONFIG_JSON > 环境变量
    this.xProjectPath = config.xProjectPath || skillsConfig.xProjectPath || process.env.SYNC_TO_GITHUB_X_PROJECT_PATH;
    this.githubRepo = config.githubRepo || skillsConfig.githubRepo || process.env.SYNC_TO_GITHUB_GITHUB_REPO;
    this.githubToken = config.githubToken || skillsConfig.githubToken || process.env.SYNC_TO_GITHUB_GITHUB_TOKEN;
    this.githubBranch = config.githubBranch || skillsConfig.githubBranch || process.env.SYNC_TO_GITHUB_GITHUB_BRANCH || 'main';
    this.defaultCategories = config.categories || skillsConfig.categories || ['技术', '博客'];
    this.defaultTags = config.tags || skillsConfig.tags || ['OpenClaw', '同步'];
    
    // 目标路径配置（默认为 _posts/ 目录，适合 Jekyll）
    this.targetPath = config.targetPath || skillsConfig.targetPath || '_posts';
    
    console.error('📋 [SyncToGitHub] 配置:');
    console.error('   xProjectPath:', this.xProjectPath || '未设置');
    console.error('   githubRepo:', this.githubRepo || '未设置');
    console.error('   githubToken:', this.githubToken ? '已设置' : '未设置');
    console.error('   githubBranch:', this.githubBranch);
    console.error('   targetPath:', this.targetPath);
    console.error('✅ [SyncToGitHub] 初始化完成');
  }

  /**
   * 从 markdown 内容中提取标题
   */
  extractTitle(content) {
    const titleMatch = content.match(/^#\s+(.+)$/m);
    if (titleMatch) {
      return titleMatch[1].trim();
    }
    const firstLines = content.split('\n').slice(0, 3).join(' ').trim();
    return firstLines.substring(0, 50) || '无标题';
  }

  /**
   * 生成文件名
   */
  generateFilename(title, date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    
    // 清理标题中的非法字符
    const cleanTitle = title
      .replace(/[<>:"/\\|？*]/g, '')
      .replace(/\s+/g, '-')
      .substring(0, 50);
    
    return `${year}-${month}-${day}-${cleanTitle}.md`;
  }

  /**
   * 转换 markdown 格式（添加 frontmatter）
   */
  async convertFormat(content, options = {}) {
    const title = options.title || this.extractTitle(content);
    const date = options.date || new Date();
    const categories = options.categories || this.defaultCategories;
    const tags = options.tags || this.defaultTags;

    const formattedDate = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:00 +0800`;

    const frontmatter = {
      title: title,
      date: formattedDate,
      categories: categories,
      tags: tags
    };

    const fileContent = matter.stringify(content, frontmatter);
    
    return {
      content: fileContent,
      filename: this.generateFilename(title, date),
      title: title
    };
  }

  /**
   * 读取源文件
   */
  async readSourceFile(sourcePath) {
    try {
      const content = await fs.readFile(sourcePath, 'utf-8');
      return content;
    } catch (error) {
      throw new Error(`读取源文件失败：${error.message}`);
    }
  }

  /**
   * 保存到本地项目
   */
  async saveToLocal(content, filename) {
    if (!this.xProjectPath) {
      throw new Error('未配置 x 项目路径 (xProjectPath)');
    }

    try {
      await fs.mkdir(this.xProjectPath, { recursive: true });
      const targetPath = path.join(this.xProjectPath, filename);
      await fs.writeFile(targetPath, content, 'utf-8');
      return targetPath;
    } catch (error) {
      throw new Error(`保存到本地项目失败：${error.message}`);
    }
  }

  /**
   * GitHub API 请求封装
   */
  async githubApiRequest(path, method = 'GET', data = null) {
    return new Promise((resolve, reject) => {
      const options = {
        hostname: 'api.github.com',
        port: 443,
        path: path,
        method: method,
        headers: {
          'Authorization': `Bearer ${this.githubToken}`,
          'Accept': 'application/vnd.github+json',
          'User-Agent': 'OpenClaw-SyncToGitHub',
          'Content-Type': 'application/json'
        }
      };

      if (data) {
        const dataString = JSON.stringify(data);
        options.headers['Content-Length'] = Buffer.byteLength(dataString);
      }

      const req = https.request(options, (res) => {
        let responseData = '';
        res.on('data', (chunk) => responseData += chunk);
        res.on('end', () => {
          try {
            const parsed = JSON.parse(responseData);
            if (res.statusCode >= 200 && res.statusCode < 300) {
              resolve(parsed);
            } else {
              reject(new Error(`GitHub API 错误 (${res.statusCode}): ${parsed.message || responseData}`));
            }
          } catch (e) {
            reject(new Error(`解析响应失败: ${responseData}`));
          }
        });
      });

      req.on('error', (error) => reject(new Error(`请求失败: ${error.message}`)));
      
      if (data) {
        req.write(JSON.stringify(data));
      }
      req.end();
    });
  }

  /**
   * 获取文件内容（用于检查文件是否存在）
   */
  async getFileContent(filePath) {
    try {
      // 对文件路径进行 URL 编码（处理中文等特殊字符）
      const encodedFilePath = encodeURIComponent(filePath).replace(/%2F/g, '/');
      const response = await this.githubApiRequest(
        `/repos/${this.githubRepo}/contents/${encodedFilePath}?ref=${this.githubBranch}`
      );
      return response;
    } catch (error) {
      if (error.message.includes('404')) {
        return null; // 文件不存在
      }
      throw error;
    }
  }

  /**
   * 创建或更新文件
   */
  async createOrUpdateFile(filePath, content, message) {
    const encodedContent = Buffer.from(content).toString('base64');
    
    // 对文件路径进行 URL 编码（处理中文等特殊字符）
    const encodedFilePath = encodeURIComponent(filePath).replace(/%2F/g, '/');
    
    // 检查文件是否已存在
    const existingFile = await this.getFileContent(filePath);
    
    const data = {
      message: message,
      content: encodedContent,
      branch: this.githubBranch
    };
    
    // 如果文件存在，需要传入 sha 来更新
    if (existingFile && existingFile.sha) {
      data.sha = existingFile.sha;
      console.error(`📝 [SyncToGitHub] 文件已存在，更新: ${filePath}`);
    } else {
      console.error(`📝 [SyncToGitHub] 创建新文件: ${filePath}`);
    }
    
    const response = await this.githubApiRequest(
      `/repos/${this.githubRepo}/contents/${encodedFilePath}`,
      'PUT',
      data
    );
    
    return response;
  }

  /**
   * 推送到 GitHub（使用 API）
   */
  async pushToGitHub(filePath, commitMessage = 'Sync markdown file') {
    if (!this.githubRepo) {
      throw new Error('未配置 GitHub 仓库 (githubRepo)');
    }

    if (!this.githubToken) {
      throw new Error('未配置 GitHub Token (githubToken)');
    }

    try {
      // 读取文件内容
      const content = await fs.readFile(filePath, 'utf-8');
      const filename = path.basename(filePath);
      
      // 目标路径
      const targetFilePath = this.targetPath 
        ? `${this.targetPath}/${filename}`
        : filename;
      
      console.error(`🚀 [SyncToGitHub] 推送到 GitHub: ${targetFilePath}`);
      
      // 使用 GitHub API 创建/更新文件
      const response = await this.createOrUpdateFile(
        targetFilePath,
        content,
        commitMessage
      );
      
      const htmlUrl = response.content?.html_url || 
                      `https://github.com/${this.githubRepo}/blob/${this.githubBranch}/${targetFilePath}`;
      
      return {
        success: true,
        message: `成功推送到 GitHub: ${this.githubRepo}`,
        pushed: true,
        branch: this.githubBranch,
        commit: commitMessage,
        filePath: targetFilePath,
        htmlUrl: htmlUrl,
        sha: response.content?.sha
      };
    } catch (error) {
      throw new Error(`推送到 GitHub 失败：${error.message}`);
    }
  }

  /**
   * 主执行函数
   */
  async execute(params = {}) {
    console.error('🚀 [SyncToGitHub] 开始执行同步任务...');
    
    const {
      sourcePath,
      content,
      title,
      date,
      categories,
      tags,
      commitMessage
    } = params;

    try {
      // 1. 获取内容
      let sourceContent = content;
      if (!sourceContent && sourcePath) {
        console.error('📖 [SyncToGitHub] 从文件读取:', sourcePath);
        sourceContent = await this.readSourceFile(sourcePath);
      }

      if (!sourceContent) {
        throw new Error('请提供源文件路径或内容');
      }

      // 2. 转换格式
      console.error('🔄 [SyncToGitHub] 转换格式...');
      const converted = await this.convertFormat(sourceContent, {
        title,
        date: date ? new Date(date) : undefined,
        categories,
        tags
      });

      console.error(`✅ [SyncToGitHub] 格式转换完成，标题："${converted.title}"`);

      // 3. 保存到本地（可选）
      let savedPath = null;
      if (this.xProjectPath) {
        console.error('💾 [SyncToGitHub] 保存到本地项目...');
        savedPath = await this.saveToLocal(converted.content, converted.filename);
        console.error(`✅ [SyncToGitHub] 已保存到本地：${savedPath}`);
      }

      // 4. 推送到 GitHub
      console.error('🚀 [SyncToGitHub] 推送到 GitHub...');
      
      // 如果有本地保存路径，从本地读取；否则直接使用转换后的内容
      const pushPath = savedPath || await this.saveToLocal(converted.content, converted.filename);
      
      const pushResult = await this.pushToGitHub(
        pushPath, 
        commitMessage || `Sync: ${converted.title}`
      );
      
      console.error(`✅ [SyncToGitHub] 推送成功！`);
      console.error(`   文件链接: ${pushResult.htmlUrl}`);

      return {
        success: true,
        message: '同步完成',
        filename: converted.filename,
        title: converted.title,
        savedPath: savedPath,
        pushResult: pushResult
      };
    } catch (error) {
      console.error('❌ [SyncToGitHub] 同步失败:', error.message);
      return {
        success: false,
        error: error.message
      };
    }
  }
}

module.exports = SyncToGitHub;
