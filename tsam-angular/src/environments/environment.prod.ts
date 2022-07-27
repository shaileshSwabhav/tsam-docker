const FILE_UPLOAD_LOACTION = "https://admin.swabhavtechlabs.com/tsm_uploads/"
// const BASE_URL = "https://swabhav-tsam.herokuapp.com/api/v1/tsam"
// const BASE_URL = "http://tsam-go-serv:8080/api/v1/tsam/"
const BASE_URL = "http://gposts.com/api/v1/tsam/"

export const environment = {
      production: true,
      BASE_URL: BASE_URL,
      // path where file is stored.
      FILE_UPLOAD_LOACTION: FILE_UPLOAD_LOACTION,

      // path to php file that handles file write on host side.
      UPLOAD_API_PATH: `${FILE_UPLOAD_LOACTION}/fileUpload.php`,

      // path with specified file name that can be used to fetch uploaded file.
      DELETE_API_PATH: `${FILE_UPLOAD_LOACTION}/fileDelete.php`,
};
