package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/jcuga/golongpoll"
)

var (
	subManager *golongpoll.LongpollManager
)

func init() {
	var err error
	subManager, err = golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled: true,
		// NOTE: if not defined here, other options have reasonable defaults,
		// so no need specifying options you don't care about
	})
	if err != nil {
		log.Fatalf("Failed to create long poll manager: %q", err)
	}

	//go generateRandomEvents(subManager)
}

func generateRandomEvents(lpManager *golongpoll.LongpollManager) {
	farm_events := []string{
		"Cow says 'Moooo!'",
		"Duck went 'Quack!'",
		"Chicken says: 'Cluck!'",
		"Goat chewed grass.",
		"Pig went 'Oink! Oink!'",
		"Horse ate hay.",
		"Tractor went: Vroom Vroom!",
		"Farmer ate bacon.",
	}
	// every 0-5 seconds, something happens at the farm:
	for {
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
		lpManager.Publish("farm", farm_events[rand.Intn(len(farm_events))])
	}
}

// Here we're providing a webpage that shows events as they happen.
func handleView(w http.ResponseWriter, r *http.Request) {
	ch := r.FormValue("ch")
	switch ch {
	case "ln":
		fmt.Fprintf(w, lnview)
	case "we":
		fallthrough
	default:
		fmt.Fprintf(w, weview)
	}
}

var weview = `
<html>
<head>
    <title>Internal Wechat Event Visualization</title>
</head>
<body>
  	<IMG SRC="/static/qr_account.jpg">
    <h1>Wechat Live Event Log</h1>
    <ul id="wechat-events"></ul>
<script src="https://code.jquery.com/jquery-1.11.3.min.js"></script>
<script>

    // for browsers that don't have console
    if(typeof window.console == 'undefined') { window.console = {log: function (msg) {} }; }

    // Start checking for any events that occurred after page load time (right now)
    // Notice how we use .getTime() to have num milliseconds since epoch in UTC
    // This is the time format the longpoll server uses.
    var sinceTime = (new Date(Date.now())).getTime();
    var category = "wxstream";

    (function poll() {
        var timeout = 45;  // in seconds
        var optionalSince = "";
        if (sinceTime) {
            optionalSince = "&since_time=" + sinceTime;
        }
        var pollUrl = "/view/events?timeout=" + timeout + "&category=" + category + optionalSince;
        // how long to wait before starting next longpoll request in each case:
        var successDelay = 10;  // 10 ms
        var errorDelay = 3000;  // 3 sec
        $.ajax({ url: pollUrl,
            success: function(data) {
                if (data && data.events && data.events.length > 0) {
                    // got events, process them
                    // NOTE: these events are in chronological order (oldest first)
                    for (var i = 0; i < data.events.length; i++) {
                        // Display event
                        var event = data.events[i];
                        $("#wechat-events").append("<li>" + event.data + " at " + (new Date(event.timestamp).toLocaleTimeString()) +  "</li>")
                        // Update sinceTime to only request events that occurred after this one.
                        sinceTime = event.timestamp;
                    }
                    // success!  start next longpoll
                    setTimeout(poll, successDelay);
                    return;
                }
                if (data && data.timeout) {
                    console.log("No events, checking again.");
                    // no events within timeout window, start another longpoll:
                    setTimeout(poll, successDelay);
                    return;
                }
                if (data && data.error) {
                    console.log("Error response: " + data.error);
                    console.log("Trying again shortly...")
                    setTimeout(poll, errorDelay);
                    return;
                }
                // We should have gotten one of the above 3 cases:
                // either nonempty event data, a timeout, or an error.
                console.log("Didn't get expected event data, try again shortly...");
                setTimeout(poll, errorDelay);
            }, dataType: "json",
        error: function (data) {
            console.log("Error in ajax request--trying again shortly...");
            setTimeout(poll, errorDelay);  // 3s
        }
        });
    })();
</script>
</body>
</html>`

var lnview = `
<html>
<head>
    <title>Internal Line Event Visualization</title>
</head>
<body>
  	<IMG SRC="https://qr-official.line.me/L/HhNE2f2q43.png">
    <h1>Line Live Event Log</h1>
    <ul id="chat-events"></ul>
<script src="https://code.jquery.com/jquery-1.11.3.min.js"></script>
<script>

    // for browsers that don't have console
    if(typeof window.console == 'undefined') { window.console = {log: function (msg) {} }; }

    // Start checking for any events that occurred after page load time (right now)
    // Notice how we use .getTime() to have num milliseconds since epoch in UTC
    // This is the time format the longpoll server uses.
    var sinceTime = (new Date(Date.now())).getTime();
    var category = "lnstream";

    (function poll() {
        var timeout = 45;  // in seconds
        var optionalSince = "";
        if (sinceTime) {
            optionalSince = "&since_time=" + sinceTime;
        }
        var pollUrl = "/view/events?timeout=" + timeout + "&category=" + category + optionalSince;
        // how long to wait before starting next longpoll request in each case:
        var successDelay = 10;  // 10 ms
        var errorDelay = 3000;  // 3 sec
        $.ajax({ url: pollUrl,
            success: function(data) {
                if (data && data.events && data.events.length > 0) {
                    // got events, process them
                    // NOTE: these events are in chronological order (oldest first)
                    for (var i = 0; i < data.events.length; i++) {
                        // Display event
                        var event = data.events[i];
                        $("#chat-events").append("<li>" + event.data + " at " + (new Date(event.timestamp).toLocaleTimeString()) +  "</li>")
                        // Update sinceTime to only request events that occurred after this one.
                        sinceTime = event.timestamp;
                    }
                    // success!  start next longpoll
                    setTimeout(poll, successDelay);
                    return;
                }
                if (data && data.timeout) {
                    console.log("No events, checking again.");
                    // no events within timeout window, start another longpoll:
                    setTimeout(poll, successDelay);
                    return;
                }
                if (data && data.error) {
                    console.log("Error response: " + data.error);
                    console.log("Trying again shortly...")
                    setTimeout(poll, errorDelay);
                    return;
                }
                // We should have gotten one of the above 3 cases:
                // either nonempty event data, a timeout, or an error.
                console.log("Didn't get expected event data, try again shortly...");
                setTimeout(poll, errorDelay);
            }, dataType: "json",
        error: function (data) {
            console.log("Error in ajax request--trying again shortly...");
            setTimeout(poll, errorDelay);  // 3s
        }
        });
    })();
</script>
</body>
</html>`
