set sample_name "Lambda";
set host_stage "false";

set sleeptime "10000";


http-get {
	set uri "/api/get";
	
	client {
		metadata {
			base64;
			header "Tmp";
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

