import { Component, ElementRef, Input, OnInit, Output, ViewChild } from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
	selector: 'app-cke-editor',
	templateUrl: './cke-editor.component.html',
	styleUrls: ['./cke-editor.component.css']
})
export class CkeEditorComponent implements OnInit {

	// Editor.
	// @ViewChild("ckeditor") ckeditor: ElementRef

	// Input's
	isImageUpload: boolean = true
	folderPath: string = "blog/images/"

	@Output() editorValue: string = ""

	// Editor configuration.
	config: any

	constructor() {
		this.editorConfig4()
	}

	ngOnInit(): void { }

	// =========================================FORMAT FUNCTIONS=================================================
	imagesFolderPath: string = this.folderPath
	imageUploadURL: string = environment.UPLOAD_API_PATH
	fileUploadLocation: string = environment.FILE_UPLOAD_LOACTION

	// Configure the editor based on the values received from the parent component.
	editorConfig4(): void {
		this.config = {
			extraPlugins: 'codeTag,kbdTag',
			removePlugins: "exportpdf",
			toolbar: [
				{ name: 'styles', items: ['Styles', 'Format'] },
				{
					name: 'basicstyles', groups: ['basicstyles', 'cleanup'],
					items: ['Bold', 'Italic', 'Underline', 'Strike', 'Subscript', 'Superscript', 'RemoveFormat', 'Code', 'Kbd']
				},
				{
					name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'],
					items: ['NumberedList', 'BulletedList', '-', 'Outdent', 'Indent', '-', 'Blockquote']
				},
				{ name: 'links', items: ['Link', 'Unlink'] },
				{ name: 'insert' }, // , items: ['SImage']
				{ name: 'document', groups: ['mode', 'document', 'doctools'], items: ['Source'] },
			],
			toolbarGroups: [
				{ name: 'styles' },
				{ name: 'basicstyles', groups: ['basicstyles', 'cleanup'] },
				{ name: 'document', groups: ['mode', 'document', 'doctools'] },
				{ name: 'paragraph', groups: ['list', 'indent', 'blocks', 'align', 'bidi'] },
				{ name: 'links' }, // Link
				{ name: 'insert' }, // Image
			],
			removeButtons: "",
			language: 'en',
			forcePasteAsPlainText: false,
			ckfinder: {
				// Upload the images to the server using the CKFinder QuickUpload command.
				uploadUrl: '/ckfinder/core/connector/php/connector.php?command=QuickUpload&type=Files&responseType=json'
			},
			// Configure your file manager integration. This example uses CKFinder 3 for PHP.
			filebrowserBrowseUrl:
				'https://ckeditor.com/apps/ckfinder/3.4.5/ckfinder.html',
			filebrowserImageBrowseUrl:
				'https://ckeditor.com/apps/ckfinder/3.4.5/ckfinder.html?type=Images',
			filebrowserUploadUrl:
				'https://ckeditor.com/apps/ckfinder/3.4.5/core/connector/php/connector.php?command=QuickUpload&type=Files',
			filebrowserImageUploadUrl:
				'https://ckeditor.com/apps/ckfinder/3.4.5/core/connector/php/connector.php?command=QuickUpload&type=Images',
			folderPath: this.imagesFolderPath,
			imageUploadURL: this.imageUploadURL,
			fileUploadLocation: this.fileUploadLocation,
			allowedContent: true,
			extraAllowedContent: 'img'
		}

		console.log("isImageUpload -> ", this.isImageUpload);


		if (this.isImageUpload) {
			this.config.toolbar.push({
				name: 'insert',
				items: ['SImage']
			})
			this.config.extraPlugins = "codeTag,kbdTag,simage"
			// this.ckConfig.toolbarGroups.push()
		}
	}

	testFunc(): void {
		console.log(" === testFunc === ");

	}

}
