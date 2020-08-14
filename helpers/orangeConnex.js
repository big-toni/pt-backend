const getItems = (sel) => {
  var items = [];
  var elements = document.body.querySelectorAll(sel);

  for (var i = 0; i < elements.length; i++) {
    var current = elements[i];
    var commonDate = current.querySelector("h3").textContent.trim()

    var entries = current.querySelectorAll("ul > div");

    entries && [...entries].map(x => {
      var item = {};
      item.date = commonDate;

      item.time = x.querySelector("span")
        .textContent
        .trim();

      item.description = x.querySelector("div[class='timeline-description']")
        .textContent
        .trim();

      let location = {}

      location.country = x.querySelector("div[class='timeline-location fl']")
      .textContent
      .replace(/\s*\n+\s*/g,'')
      .trim();

      item.location = location

      items.push(item)
    })
  }
  return items
};