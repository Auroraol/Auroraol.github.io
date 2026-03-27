/**
 * Set default theme to dark on first visit
 */
export function setDefaultTheme() {
  // Check if user has already set a theme preference
  const savedMode = sessionStorage.getItem('mode');
  
  // If no preference saved and no system preference set, default to dark
  if (!savedMode) {
    const html = document.documentElement;
    const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // Set dark theme as default
    html.setAttribute('data-mode', 'dark');
    sessionStorage.setItem('mode', 'dark');
    
    // Dispatch event to notify other components
    window.dispatchEvent(new CustomEvent('theme-change', {
      detail: { mode: 'dark' }
    }));
  }
}