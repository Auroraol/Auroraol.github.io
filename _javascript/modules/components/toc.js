function directChildByTag(element, tagName) {
  return [...element.children].find((child) => child.tagName === tagName);
}

function initTocCollapseControls() {
  const tocRoot = document.getElementById('toc');

  if (!tocRoot) {
    return;
  }

  tocRoot.querySelectorAll('li').forEach((item, index) => {
    const link = directChildByTag(item, 'A');
    const childList =
      directChildByTag(item, 'UL') || directChildByTag(item, 'OL');

    if (!link || !childList) {
      return;
    }

    item.classList.add('toc-collapsible');

    if (!childList.id) {
      childList.id = `toc-sublist-${index}`;
    }

    const button = document.createElement('button');
    button.type = 'button';
    button.className = 'toc-toggle';
    button.setAttribute('aria-expanded', 'true');
    button.setAttribute('aria-controls', childList.id);
    button.setAttribute('aria-label', `折叠 ${link.textContent.trim()}`);

    button.addEventListener('click', (event) => {
      event.preventDefault();
      event.stopPropagation();

      const collapsed = item.classList.toggle('toc-collapsed');
      button.setAttribute('aria-expanded', String(!collapsed));
      button.setAttribute(
        'aria-label',
        `${collapsed ? '展开' : '折叠'} ${link.textContent.trim()}`
      );
    });

    item.insertBefore(button, link);
  });
}

export function toc() {
  if (document.querySelector('.content h1, .content h2, .content h3')) {
    // see: https://github.com/tscanlin/tocbot#usage
    tocbot.init({
      tocSelector: '#toc',
      contentSelector: '.content',
      ignoreSelector: '[data-toc-skip]',
      headingSelector: 'h1, h2, h3, h4',
      collapseDepth: 6,
      orderedList: false,
      scrollSmooth: false
    });

    initTocCollapseControls();
  }
}
