;(function() {
    var $passwordModal = $('#passwordModal');

    $('#password_form').submit(function(e) {
        var _this = $(this);
        var $btns = _this.find('button');

        var url = _this.attr('action');
        var data = _this.serialize();

        $btns.attr('disabled', 'true');

        $.post(url, data, function(json) {
            if(json.success){
                iSuccess(json.msg, function() {
                    $passwordModal.modal('hide');
                });
            }else{
                $btns.removeAttr('disabled');
                iError(json.msg);
            }
        }, 'json');

        return false;
    });
})();