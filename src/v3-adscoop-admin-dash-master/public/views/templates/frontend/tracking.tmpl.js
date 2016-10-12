(function() {

	var asStartTime = (new Date()).getTime();

	var asOldOnBeforeUnload = null;

	if (typeof window.onbeforeunload === 'function') {
		asOldOnBeforeUnload = window.onbeforeunload;
	}

	window._ast = (function() {
		var events = [];
		var dome = {
			get: function() {
				return events;
			},
			fireEvent: function(pingData) {
				events.push(pingData)
			}
		}
		return dome;
	})(window);

	function getParameterByName(name) {
	    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
	    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
	        results = regex.exec(location.search);
	    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
	}

	var atf = getParameterByName('_ast').split("_")[0];

	var callback = 'ast' + (Date.now() / 100 | 0)

	var trck = document.createElement('script');
	trck.type = 'text/javascript'; trck.async = true;
	trck.src = 'http://{{ .Host }}/loadClient?_ast=' + getParameterByName('_ast') + '&callback=' + callback;
	var s = document.getElementsByTagName('script')[0];
	if (getParameterByName('_ast') != "") {
		s.parentNode.insertBefore(trck, s);
	}

	window[callback] = function(data){
		if (typeof loadAfterAst != 'undefined') {
			loadAfterAst();
		}
		tempEngagements = _ast.get();

		window._ast = (function() {

			if (data.ET) {
				var timeout = Math.floor(Math.random() * (data.MaxTimeout-data.MinTimeout)) + data.MinTimeout;

				setTimeout(function() {
					window.location = data.Redir;
				}, timeout * 1000)
			}

		if (data.EnableUnloadTracking === true) {


			window.onbeforeunload = function() {
				if (typeof asOldOnBeforeUnload === 'function') {
					asOldOnBeforeUnload();
				}

				asStopTime = (new Date()).getTime();

				timeOnPage = ((asStopTime - asStartTime) / 1000)

				dome.fireUnload(timeOnPage);
			}
		}

	if (typeof window.astTimeout == 'undefined') {
		window.astTimeout = 1;
	} else {
		window.astTimeout = window.astTimeout * 1000;
	}

	var atf = getParameterByName('_ast').split("_")[0];

	var bbsiSet = false;
	var bbsiLoaded = false;
	var atf_trackLoad = false;
	var atf_track = false;

	if (typeof getParameterByName('_ast').split("_")[1] != "undefined" && getParameterByName('_ast').split("_")[1] == "1") {
		bbsiSet = true

		var s = document.createElement('script');
		s.src = "//a.bbsi.io/e4faecc5a7cf7a3d2e2db72292bed08f/a.js";
		document.head.appendChild(s);

		s.onload = function() {
			bbsiLoaded = true;

			if (atf_track) {
				bbsiBA.dingDong("engagement_track");
			}

			if (atf_trackLoad) {
				bbsiBA.dingDong("engagement_trackLoad");
			}
		}
	}
	function getParameterByName(name) {
	    name = name.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
	    var regex = new RegExp("[\\?&]" + name + "=([^&#]*)"),
	        results = regex.exec(location.search);
	    return results === null ? "" : decodeURIComponent(results[1].replace(/\+/g, " "));
	}
	var dome = {
		get: function(selector) {

		},
		fireEvent: function(pingData) {
			if (typeof pingData == "undefined") {
				pingData = {};
			}
			if (typeof pingData.player_instance_id == 'undefined') {
				pingData.player_instance_id = "";
				pingData.tos = "";
			}
			atf_track = true;
			if (!atf) {
				return;
			}
			dome.track(data.Host, pingData);

			if (bbsiLoaded) {
				bbsiBA.dingDong("engagement_track");
			}

		},
		fireUnload: function(tos) {
			pingData = {};
			pingData.player_instance_id = "";
			pingData.tos = "";

			pingData.tos = tos;

			console.log("tos", tos);

			dome.trackUnload(data.Host, pingData)
		},
		track: function(host, pingData) {
			dome.ping(host, 'engagement', pingData);
		},
		trackUnload: function(host, pingData) {
			dome.ping(host, 'tos', pingData);
		},
		trackLoad: function(host) {
			var pingData = {};
			pingData.player_instance_id = ""
			pingData.tos = ""
			dome.ping(host, 'load', pingData);
		},
		ping: function(host, uri, pingData) {
			var img = document.createElement("img");
			img.style.width = "1px";
			img.style.opacity = "0";
			img.style.height = "1px";
			img.setAttribute("src", host + "/t/" + uri + "?_cb=" + Math.random()+ "&player_instance_id=" + pingData.player_instance_id + "&_tos=" + pingData.tos)
			document.body.appendChild(img);
		}
	}
	var oldonload = window.onload;
	window.onload = function() {
		if (typeof oldonload === 'function') {
				oldonload();
		}
		setTimeout(function() {

			atf_trackLoad = true;
			dome.trackLoad(data.Host);

			if (bbsiLoaded) {
				bbsiBA.dingDong("engagement_trackLoad");
			}

		}, window.astTimeout);
	}

	return dome
})(window, data);
	for (var i = 0; i < tempEngagements.length; i++) {
		_ast.fireEvent(tempEngagements[i]);
	}
}
})(window);
