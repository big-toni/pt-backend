const getItems = (sel) => {
  var items = [];
  var entries = document.body.querySelectorAll("div[class=el-collapse-item__content] > div > div[data-v-41aef011][class='detail-list col-md-12']");

  entries && [...entries].map(x => {
    var item = {};

    var infos = x.querySelectorAll("p[data-v-41aef011]")

    item.time = infos[0]
      .textContent
      .trim();

    item.description = infos[1]
      .textContent
      .trim();

    let location = {}

    location.country = infos[2]
    .textContent
    .replace(/\s*\n+\s*/g,'')
    .trim();

    item.location = location

    items.push(item)
  })
  return items
};