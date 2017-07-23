;(function() {
    var $body = $('body');

    $('#ajax_form').submit(function(e) {
        var _this = $(this);
        var $btns = _this.find('button');

        var url = _this.attr('action');
        var data = _this.serialize();

        $btns.attr('disabled', 'true');

        $.post(url, data, function(json) {
            if(json.success){
                iSuccess(json.msg, function() {
                    if (json.redirect) {
                        location.href = json.redirect;
                    }
                });
            }else{
                $btns.removeAttr('disabled');

                iError(json.msg, function() {
                    if (json.redirect) {
                        location.href = json.redirect;
                    }
                });
            }
        }, 'json');

        return false;
    });

    $body.on('submit', '.ajax_form',function(e) {
        var _this = $(this);
        var $btns = _this.find('button');

        var url = _this.attr('action');
        var data = _this.serialize();

        $btns.attr('disabled', 'true');

        $.post(url, data, function(json) {
            if(json.success){
                iSuccess(json.msg, function() {
                    if (json.redirect) {
                        location.href = json.redirect;
                    }
                });
            }else{
                $btns.removeAttr('disabled');

                iError(json.msg, function() {
                    if (json.redirect) {
                        location.href = json.redirect;
                    }
                });
            }
        }, 'json');

        return false;
    });
})();