import re
import pandas as pd
import matplotlib.pyplot as plt
import json

# Function to parse the almost-JSON data from the file content
def parse_data(text_data):
    objects = re.findall(r'\{.*?\}', text_data, re.DOTALL)
    if not objects:
        return None
    json_str = f'[{','.join(objects)}]'
    json_str = re.sub(r'NumberLong\((\"?)(\d+)(\1)\)', r'\2', json_str)
    try:
        data = json.loads(json_str)
        return data
    except json.JSONDecodeError as e:
        print(f"Error decoding JSON: {e}")
        return None

# Read data from file
try:
    with open('5.txt', 'r', encoding='utf-8') as f:
        content = f.read()
except FileNotFoundError:
    print("错误：未找到 '新建 文本文档.txt' 文件。")
    exit()

# Parse the content of the file
parsed_data = parse_data(content)

if parsed_data:
    durations = [item.get('res_isv_timestamp') for item in parsed_data if 'res_isv_timestamp' in item]
    
    if not durations:
        print("未能在文件中找到 'res_isv_timestamp' 数据。")
    else:
        df = pd.Series(durations)
        value_counts = df.value_counts().sort_index()

        # --- Dynamic Height Calculation ---
        # Calculate figure height based on the number of unique data points.
        # Use a base height and add height for each item. Min height is 8 inches.
        num_items = len(value_counts)
        fig_height = max(8, num_items * 0.4) # 0.4 inches per bar

        try:
            plt.rcParams['font.sans-serif'] = ['SimHei']
            plt.rcParams['axes.unicode_minus'] = False
        except Exception as e:
            print("警告：设置中文字体失败。图表中的中文可能无法正常显示。")

        # Plotting the horizontal bar chart with dynamic height
        plt.figure(figsize=(10, fig_height))
        value_counts.plot(kind='barh')
        
        plt.title('耗时计数图 (Duration Count)')
        plt.xlabel('数量 (Count)')
        plt.ylabel('耗时 (s)')
        plt.grid(axis='x', alpha=0.75)

        output_filename = 'duration_count_bar_chart_5.png'
        plt.savefig(output_filename, bbox_inches='tight')

        print(f"已生成动态高度的水平条形图: {output_filename}")
else:
    print("未能在文件中解析出任何数据块。")
