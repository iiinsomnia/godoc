;(function() {
	// 搜索
	var $searchForm = $('#search_form');
	var $searchSubmit = $('#search_submit');
	// 数据
	var $dataBody = $('#data_body');
    var $count = $('#count');
	// 分页
	var $pagination = $('#ajax_pagination');
	var $firstPage = $('#first_page');
	var $prevPage = $('#prev_page');
	var $nextPage = $('#next_page');
	var $lastPage = $('#last_page');
	var $skipBtn = $('#skip_btn');
	var $pageNum = $('#page_num');
    var $curPage = $('#cur_page');
    var $pages = $('#pages');

    var curPage = 1;
    var lastPage = parseInt($pages.data('num'));
    var pageUrl = $dataBody.data('url');
    var query = {};

	// 提交搜索
	$searchSubmit.click(function(e) {
		/* Act on the event */
		if ($searchForm) {
			formData = $searchForm.serializeArray();

			formData.forEach(function(ele, index, arr) {
				query[ele.name] = ele.value;
			});
		}

		query["page"] = 1;

		$.get(pageUrl, query, function(json) {
			if(json.success) {
				if(json.data.pages > 1) {
					$firstPage.attr('disabled', 'true');
					$prevPage.attr('disabled', 'true');
					$nextPage.removeAttr('disabled');
					$lastPage.removeAttr('disabled');
					$pagination.show();
				} else {
					$pagination.hide();
				}

				curPage = 1;
				lastPage = json.data.pages;

				$count.text(json.data.count);
				$curPage.text(1);
				$pages.text(json.data.pages);

				$dataBody.html(json.data.html);
			} else {
				iError(json.msg, function() {
                    if (json.redirect) {
                        location.href = json.redirect;
                    }
                });
			}
		}, 'json');
	});

	// 首页
	$firstPage.click(function(e) {
		/* Act on the event */
		if(curPage > 1){
			curPage = 1;
			$curPage.text(curPage);

			query["page"] = curPage;

			$.get(pageUrl, query, function(json) {
				if(json.success) {
					$firstPage.attr('disabled', 'true');
					$prevPage.attr('disabled', 'true');
					$nextPage.removeAttr('disabled');
					$lastPage.removeAttr('disabled');

					$dataBody.html(json.data.html);
				} else {
					iError(json.msg, function() {
						if (json.redirect) {
							location.href = json.redirect;
						}
					});
				}
			}, 'json');
		}

		return false;
	});

	// 前一页
	$prevPage.click(function(e) {
		/* Act on the event */
		if(curPage > 1){
			curPage--;
			$curPage.text(curPage);

			query["page"] = curPage;

			$.get(pageUrl, query, function(json) {
				if (json.success) {
					if(curPage == 1) {
						$firstPage.attr('disabled', 'true');
						$prevPage.attr('disabled', 'true');
					}
					$nextPage.removeAttr('disabled');
					$lastPage.removeAttr('disabled');

					$dataBody.html(json.data.html);
				} else {
					iError(json.msg, function() {
						if (json.redirect) {
							location.href = json.redirect;
						}
					});
				}
			}, 'json');
		}

		return false;
	});

	// 后一页
	$nextPage.click(function(e) {
		/* Act on the event */
		if(curPage < lastPage){
			curPage++;
			$curPage.text(curPage);

			query["page"] = curPage;

			$.get(pageUrl, query, function(json) {
				if(json.success){
					$firstPage.removeAttr('disabled');
					$prevPage.removeAttr('disabled');

					if(curPage == lastPage){
						$nextPage.attr('disabled', 'true');
						$lastPage.attr('disabled', 'true');
					}

					$dataBody.html(json.data.html);
				} else {
					iError(json.msg, function() {
						if (json.redirect) {
							location.href = json.redirect;
						}
					});
				}
			}, 'json');
		}

		return false;
	});

	// 尾页
	$lastPage.click(function(e) {
		/* Act on the event */
		if(curPage < lastPage){
			curPage = lastPage;
			$curPage.text(curPage);

			query["page"] = curPage;

			$.get(pageUrl, query, function(json) {
				if(json.success) {
					$firstPage.removeAttr('disabled');
					$prevPage.removeAttr('disabled');
					$nextPage.attr('disabled', 'true');
					$lastPage.attr('disabled', 'true');

					$dataBody.html(json.data.html);
				} else {
					iError(json.msg, function() {
						if (json.redirect) {
							location.href = json.redirect;
						}
					});
				}
			}, 'json');
		}

		return false;
	});

	// 跳转页
	$skipBtn.click(function(e) {
		/* Act on the event */
		var pageNum = $pageNum.val();

		if(pageNum.trim() == '' || isNaN(pageNum)) {
			iError('请输入页码');
			return false;
		}

		if(pageNum < 1) {
			iError('页码太小了');
			return false;
		}

		if(pageNum > lastPage) {
			iError('页码太大了');
			return false;
		}

		curPage = pageNum;
		$curPage.text(curPage);

		query["page"] = curPage;

		$.get(pageUrl, query, function(json) {
			if(json.success){
				if(curPage == 1) {
					$firstPage.attr('disabled', 'true');
					$prevPage.attr('disabled', 'true');
					$nextPage.removeAttr('disabled');
					$lastPage.removeAttr('disabled');
				} else if(curPage == lastPage) {
					$firstPage.removeAttr('disabled');
					$prevPage.removeAttr('disabled');
					$nextPage.attr('disabled', 'true');
					$lastPage.attr('disabled', 'true');
				} else {
					$firstPage.removeAttr('disabled');
					$prevPage.removeAttr('disabled');
					$nextPage.removeAttr('disabled');
					$lastPage.removeAttr('disabled');
				}

				$dataBody.html(json.data.html);
			}else{
				iError(json.msg, function() {
                    if (json.redirect) {
                        location.href = json.redirect;
                    }
                });
			}
		}, 'json');
	});
})();