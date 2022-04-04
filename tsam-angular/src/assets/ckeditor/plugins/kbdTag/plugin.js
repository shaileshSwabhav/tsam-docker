CKEDITOR.plugins.add( 'kbdTag', {
  icons: 'kbd',
  onLoad: function() {
    // Adds CSS 
    CKEDITOR.addCss(
      'kbd {' +
        'background-color: #000;' +
        'color: #fff;' +
      '}'
    )
  },
  init: function( editor ) {
    // editor.addContentsCss( this.path + 'styles.css' );
    editor.addCommand( 'wrapKbd', {
      exec: function( editor ) {
        editor.insertHtml( '<kbd>' + editor.getSelection().getSelectedText() + '</kbd>' );
      }
    });
    editor.ui.addButton( 'Kbd', {
      label: 'Wrap kbd',
      command: 'wrapKbd',
      toolbar: 'insert'
    });
  }
});