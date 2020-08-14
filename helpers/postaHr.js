const getItems = (sel) => {
  let items = [];
  let elements = document.body.querySelectorAll(sel);

  for (var i = 0; i < elements.length; i++) {
    let current = elements[i];

    let blocks = current.querySelectorAll("div[class='styles__block___2UkH9']");

    blocks && [...blocks].map(bl => {

      let entries = bl.querySelectorAll("div[class='styles__row___UT_E9']");
      var commonStatus = null
      entries && [...entries].map(x => {
        let item = {};

        const commonStatusElement = x.querySelector("div[class='styles__category___2lOSC']")
        if (commonStatusElement) {
          commonStatus = commonStatusElement.textContent.trim();
          return
        }

        item.status = commonStatus;

        const dateElement = x.querySelector("div[class='styles__date___-QZuh']")
          .textContent
          .trim();

        const timeRE = /[0-9]+\:[0-9]*\:[0-9]*\s[A-Z]+/g;
        const timeMatch = dateElement.match(timeRE);
        item.time = timeMatch && timeMatch[0];

        const dateRE = /[0-9]+\/[0-9]+\/[0-9]+/g;
        const dateMatch = dateElement.match(dateRE);
        item.date = dateMatch && dateMatch[0];

        const infoElement = x.querySelector("div[class='styles__status___3MLpj']")
        const infoStr = infoElement && infoElement.textContent.trim()

        const infoData = infoStr && infoStr.split("-")

        if (infoData) {
          const locStr = infoData[0] && infoData[0].trim();

          if (locStr) {
            let location = {}
            const zipRE = /[0-9]+/g;
            const zipMatch = locStr.match(zipRE);
            location.zip = zipMatch && zipMatch[0];
  
            const cityRE = /[A-Za-z]+/g;
            const cityMatch = locStr.match(cityRE);
            location.city = cityMatch && cityMatch[0].trim();
  
            item.location = location;
          }
          
          item.description = infoData[1] && infoData[1].trim();
        }

        items.push(item);
      });
    });
  }
  return items;
};