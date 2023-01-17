set sample_name "Lambda";
set host_stage "false";

set sleeptime "10000";
set useragent "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:47.0) Gecko/20100101 Firefox/47.0";

http-get {
	set uri "/api/get";
	
	client {
		metadata {
			base64;
			header "Ping";
		}
	}

	server {

		output {
            base64;
			print;
		}
	}
}

http-post {
	set uri "/api/post";

	client {
		id {
			append "/default.asp";
			uri-append;
		}

		header "Pong";

		output {
            base64;
			print;
		}
	}

	server {
		output { 
            base64;
			print;
		}
	}
}

