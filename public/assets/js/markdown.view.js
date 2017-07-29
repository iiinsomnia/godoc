;(function () {
    editorView = editormd.markdownToHTML("editormd_view", {
        htmlDecode: "style,script,iframe",
        taskList: true
    });
})();