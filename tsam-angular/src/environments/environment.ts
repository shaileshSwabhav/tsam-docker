// This file can be replaced during build by using the `fileReplacements` array.
// `ng build --prod` replaces `environment.ts` with `environment.prod.ts`.
// The list of file replacements can be found in `angular.json`.

const FILE_UPLOAD_LOACTION = "https://admin.swabhavtechlabs.com/tsm_uploads/testing/"
const BASE_URL = "http://127.0.0.1:8080/api/v1/tsam"

export const environment = {
  production: false,
  BASE_URL: BASE_URL,
  // path where file is stored.
  FILE_UPLOAD_LOACTION: FILE_UPLOAD_LOACTION,

  // path to php file that handles file write on host side.
  UPLOAD_API_PATH: `${FILE_UPLOAD_LOACTION}/fileUpload.php`,

  // path with specified file name that can be used to fetch uploaded file.
  DELETE_API_PATH: `${FILE_UPLOAD_LOACTION}/fileDelete.php`,
};

/*
 * For easier debugging in development mode, you can import the following file
 * to ignore zone related error stack frames such as `zone.run`, `zoneDelegate.invokeTask`.
 *
 * This import should be commented out in production mode because it will have a negative impact
 * on performance if an error is thrown.
 */
// import 'zone.js/dist/zone-error';  // Included with Angular CLI.
