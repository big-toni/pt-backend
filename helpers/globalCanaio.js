const getItems = (sel) => {
  var items = [];
  var entries = document.body.querySelectorAll(
    "div[class='TrackingDetail--step--2sEzXUy']"
  );

  // return entries;

  entries &&
    [...entries].map((x) => {
      var item = {};

      // var infos = x.querySelectorAll(
      //   "div[class='TrackingDetail--stepContent--2n0dtpG undefined TrackingDetail--stepContentHeadGrey--13kDuUL']"
      // );

      const dateElement = x
        .querySelector("div[class='TrackingDetail--timeInfoWrap--Ad4suAI']")
        .textContent.trim();

      item.time = dateElement;

      const descriptionElement = x
        .querySelector("[class='TrackingDetail--head--20GpNSP']")
        .textContent.trim();

      item.description = descriptionElement;

      // let location = {};

      // location.country = infos[2].textContent.replace(/\s*\n+\s*/g, '').trim();

      // item.location = location;

      items.push(item);
    });
  return items;
};

// TrackingDetail--timeInfoWrap--Ad4suAI
