;(function () {
    var $body = $('body');

    $body.on('click', '.reset', function(e) {
        var _this = $(this);

        if (confirm('确定要重置密码？')) {
            var url = _this.data('url');

            $.get(url, function (json) {
                if(json.success){
                    iSuccess(json.msg);
                }else{
                    iError(json.msg);
                }
            }, 'json');
        }
    });

    $body.on('click', '.delete', function(e) {
        var _this = $(this);

        var msg = _this.data('msg');

        if (!msg) {
            msg = '确定要删除？'
        }

        if (confirm(msg)) {
            var url = _this.data('url');

            $.get(url, function (json) {
                if(json.success){
                    iSuccess(json.msg, function() {
                        if (json.redirect) {
                            location.href = json.redirect;
                        }
                    });
                }else{
                    iError(json.msg);
                }
            }, 'json');
        }
    });
})();