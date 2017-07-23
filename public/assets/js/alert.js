;(function () {
    toastr.options = {
        "timeOut": 1000
    };

    function iSuccess(msg, callback) {
        if (msg instanceof(Array)) {
            msg = msg.join("<br/>");
        }

        toastr.success(msg);

        if(callback){
            setTimeout(callback, 1000);
        }
    }

    function iInfo(msg, callback) {
        if (msg instanceof(Array)) {
            msg = msg.join("<br/>");
        }

        toastr.info(msg);

        if(callback){
            setTimeout(callback, 1000);
        }
    }

    function iWarning(msg, callback) {
        if (msg instanceof(Array)) {
            msg = msg.join("<br/>");
        }

        toastr.warning(msg);

        if(callback){
            setTimeout(callback, 1000);
        }
    }

    function iError(msg, callback) {
        if (msg instanceof(Array)) {
            msg = msg.join("<br/>");
        }

        toastr.error(msg);

        if(callback){
            setTimeout(callback, 1000);
        }
    }

    window.iSuccess = iSuccess;
    window.iInfo = iInfo;
    window.iWarning = iWarning;
    window.iError = iError;
})();