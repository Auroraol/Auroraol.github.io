/**
 * Expand or close the sidebar in mobile screens.
 */

const ATTR_DISPLAY = 'sidebar-display';
const STORAGE_KEY = 'sidebar-collapsed';

class SidebarUtil {
  static isExpanded = false;

  static toggle() {
    if (SidebarUtil.isExpanded === false) {
      document.body.setAttribute(ATTR_DISPLAY, '');
    } else {
      document.body.removeAttribute(ATTR_DISPLAY);
    }

    SidebarUtil.isExpanded = !SidebarUtil.isExpanded;
  }

  // Desktop collapse functionality
  static initDesktopCollapse() {
    const collapseBtn = document.getElementById('sidebar-collapse');
    if (!collapseBtn) return;

    // Restore saved state
    const isCollapsed = localStorage.getItem(STORAGE_KEY) === 'true';
    if (isCollapsed) {
      document.body.setAttribute('data-sidebar-collapsed', 'true');
    }

    // Toggle collapse on button click
    collapseBtn.addEventListener('click', () => {
      const currentState = document.body.getAttribute('data-sidebar-collapsed');
      const newState = currentState !== 'true';
      
      if (newState) {
        document.body.setAttribute('data-sidebar-collapsed', 'true');
      } else {
        document.body.removeAttribute('data-sidebar-collapsed');
      }
      
      // Save state to localStorage
      localStorage.setItem(STORAGE_KEY, newState);
    });
  }
}

export function sidebarExpand() {
  document
    .getElementById('sidebar-trigger')
    .addEventListener('click', SidebarUtil.toggle);

  document.getElementById('mask').addEventListener('click', SidebarUtil.toggle);
  
  // Initialize desktop collapse
  SidebarUtil.initDesktopCollapse();
}
