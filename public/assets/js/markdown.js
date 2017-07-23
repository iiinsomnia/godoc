;(function () {
    editor = editormd("editormd", {
        width: "100%",
        height: 600,
        path: "/assets/vendor/editor.md/lib/",
        toolbarIcons: function() {
            return [
                "undo", "redo", "|",
                "bold", "del", "italic", "quote",  "|",
                "h1", "h2", "h3", "h4", "h5", "h6", "|",
                "list-ul", "list-ol", "hr", "|",
                "table", "link", "datetime", "clear", "|",
                "watch", "preview", "info"
            ]
        }
    });
})();