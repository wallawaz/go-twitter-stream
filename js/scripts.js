function updatePooScore() {
  let counter = document.querySelector('.score');
  counter.innerHTML = String(POOCOUNTER);
}

function updateTweets(tweet) {
    let img = ' <img src="' +
         tweet["user"]["profile_image_url_https"] +
         '" alt="" style="width:50px; height:auto;">';

    let dt = $("#mainTable").DataTable();
    //console.log(tweet);

    TWEETCOUNTER+=1;
    let row = [
      TWEETCOUNTER,
      tweet["created_at"],
      img,
      tweet["user"]["screen_name"],
      tweet["text"]
    ];
    dt.row.add(row);
    dt.order([0, "desc"]).draw();

    if (TWEETCOUNTER > NUMROWS) {
      const tableIDs = dt.rows()[0]
      const removeID = tableIDs[tableIDs.length - 1];
      dt.row(removeID).remove();
      dt.draw()
    }
}

function pushTweet(tweetID) {
  let poopDiv = document.createElement("div");
  poopDiv.className = "poopDiv";
  poopDiv.id = "poo-" + POOCOUNTER;

  //get these randomly
  //poopDiv.style.position = "absolute";

  let pos = randomPosition();

  poopDiv.style.top = pos.y;
  poopDiv.style.left = pos.x;
  main.appendChild(poopDiv);
  console.dir(poopDiv);


  let poopImage = document.createElement("img");

  poopImage.src = imageURLs[Math.round(Math.random(0,1))]
  poopImage.alt = "poo-" + POOCOUNTER;

  poopImage.style.width = '50px';
  poopImage.style.height = '50px';
  poopDiv.appendChild(poopImage);


  //poopDiv.appendChild(tweetSpan);
  twttr.widgets.createTweet(
      tweetID,
      poopDiv,
      {
          theme: 'dark'
      }
  );
  POOCOUNTER += 1;
  updatePooScore();
  }

function randomPosition() {
    const pos = getWindowPx();
    console.dir(pos);

    const x = Math.floor(Math.random() * pos.x) + 1;
    const y = Math.floor(Math.random() * pos.y) + 1;
    return {"x": `${x}px`, "y": `${y}px`}
}
