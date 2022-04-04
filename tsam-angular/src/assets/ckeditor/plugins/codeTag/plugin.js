CKEDITOR.plugins.add( 'codeTag', {
  icons: 'code',
  onLoad: function() {
    // Adds CSS 
    CKEDITOR.addCss(
      'code {' +
        'background-color: #ffeaee;' +
        'color: #fa002a;' +
      '}'
    )
  },
  init: function( editor ) {
    // editor.addContentsCss( this.path + 'styles.css' );
    editor.addCommand( 'wrapCode', {
      exec: function( editor ) {
        editor.insertHtml( '<code>' + editor.getSelection().getSelectedText() + '</code>' );
      }
    });
    editor.ui.addButton( 'Code', {
      label: 'Wrap code',
      command: 'wrapCode',
      toolbar: 'insert'
    });
  }
});