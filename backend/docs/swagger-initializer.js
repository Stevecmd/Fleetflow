
/**
 * window.onload handler to initialize swagger-ui.
 *
 * This function will be called when the page has finished loading.
 * It will initialize swagger-ui and load the swagger.json file
 * at the given url.
 *
 * @param {string} url - The url of the swagger.json file.
 * @param {string} dom_id - The id of the dom element where the swagger-ui
 *   should be rendered.
 * @param {boolean} deepLinking - If true, enable deep linking for the swagger-ui.
 * @param {array} presets - An array of presets to apply to the swagger-ui.
 * @param {array} plugins - An array of plugins to apply to the swagger-ui.
 * @param {string} layout - The layout to use for the swagger-ui.
 */
window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "swagger.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};
