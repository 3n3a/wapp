import "htmx.org"
import './lib_htmx'

import "morphdom"
import "htmx.org/dist/ext/morphdom-swap"

/* All Fields Filled */
window.allFilled = function (formSelector) {
  // Check if all form fields are filled
  const form = document.querySelector(formSelector);
  const inputs = form?.querySelectorAll("input[required]");
  if (!inputs) return false

  for (const input of inputs) {
    if (!input.value.trim()) {
      return false;
    }
  }

  return true;
}

// Show Element by Selector
window.showEl = (elSelector) => {
  htmx.find(elSelector).style.display = 'block'
}

// Remove Elements by Selector
window.cleanEls = (elSelector) => {
  htmx.findAll(elSelector).forEach(e => e.remove())
}

window.hideLoadingIndicators = (selector = '#loading-indicator') => {
  htmx.findAll(selector).forEach(e => {
    e.style.display = 'none'
  })
}


htmx.onLoad(() => {
  // Hide Loading Indicator if Empty Form
  // if (!allFilledSubdomainsForm() || !allFilledDnsForm()) {
    window.hideLoadingIndicators()
  // }
})
// htmx.logAll();
