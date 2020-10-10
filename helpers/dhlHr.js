const getItems = (sel) => {
  var item = {};

  let element = document.body.querySelector("div[class='l-grid l-grid--w-100pc-s l-grid--w-auto-m']")
  const infoElement = element.querySelector("h3[class='c-tracking-result--status-copy-message']")
    .textContent
    .trim();

  const dateElement = element.querySelector("div[class='c-tracking-result--status-copy-date  ']")
    .textContent
    .trim();

  item.description = infoElement;
  item.time = dateElement.split(' ').slice(0, 4).join(' ');

  return [item];
};