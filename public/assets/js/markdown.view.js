;(function () {
    editorView = editormd.markdownToHTML("editormd_view", {
        htmlDecode: "style,script,iframe",
        emoji: true,
        taskList: true,
    });
})();