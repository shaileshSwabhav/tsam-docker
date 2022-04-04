import { Injectable } from '@angular/core';
import * as CryptoJS from 'crypto-js';

const SECURE_STORAGE = require('secure-web-storage');
const SECRET_KEY = 'secret_key';

@Injectable({
	providedIn: 'root'
})
export class StorageService {

	constructor() { }

	// Set Data To Local Storage
	setItem(key: string, data: any) {
		sessionStorage.setItem(key, data);
	}


	// Get Item By Key.
	getItem(key: string) {
		return sessionStorage.getItem(key);
	}

	// Remove Item By Key.
	removeItem(key: string) {
		return sessionStorage.removeItem(key);
	}

	public secureStorage = new SECURE_STORAGE(localStorage, {
		// Encrypt the localstorage data
		hash: function hash(key) {
			key = CryptoJS.SHA256(key, SECRET_KEY);
			return key.toString();
		},
		encrypt: function encrypt(data: any) {
			data = CryptoJS.AES.encrypt(data, SECRET_KEY);
			data = data.toString();
			return data;
		},
		// Decrypt the encrypted data
		decrypt: function decrypt(data: any) {
			data = CryptoJS.AES.decrypt(data, SECRET_KEY);
			data = data.toString(CryptoJS.enc.Utf8);
			return data;
		}
	});
}
