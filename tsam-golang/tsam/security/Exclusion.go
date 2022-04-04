package security

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/techlabs/swabhav/tsam/errors"
	"github.com/techlabs/swabhav/tsam/web"
)

// RegisterRoutes will exclude routes on which middleware should not be applied
func (auth *Authentication) RegisterRoutes(excludedRoutes []*mux.Route) func(http.Handler) http.Handler {
	// for loop creates the regexp array based on the route
	// eg: /college -> ^/college$
	// eg: /college/{collegeID} -> ^/college/(?P<v0>[^/]+)$ (regex is example and not literal conversion)
	// eg: /api/v1/tsam/login -> ^/api/v1/tsam/login[/]?$ (literal conversion)
	var excludedRoutesRegexp []*regexp.Regexp
	rl := len(excludedRoutes)
	for i := 0; i < rl; i++ {
		r := excludedRoutes[i]
		pathRegexp, _ := r.GetPathRegexp()
		regx, _ := regexp.Compile(pathRegexp)
		excludedRoutesRegexp = append(excludedRoutesRegexp, regx)
	}

	// fmt.Println("=============================================")
	// fmt.Println("Exclude Routes -> ", excludedRoutesRegexp)
	// fmt.Println("=============================================")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			exclude := false

			// requestMethod will return the method type such as GET, POST, PUT, DELETE
			requestMethod := r.Method
			// for loop checks through all the routes in the excludedRoutes slice
			// if we find a route which is excluded then its uri form is created and redirected to that page
			// if not found the middleware is applied to it
			for i := 0; i < rl; i++ {
				excludedRoute := excludedRoutes[i]
				// methods will store the method type such as GET
				methods, _ := excludedRoute.GetMethods()
				ml := len(methods)

				// methodMatched will be used to check if the given method is excluded or not
				methodMatched := false
				if ml < 1 {
					methodMatched = true
				} else {
					//  checks if the method matches the requestMethod type
					// i.e. if requestMethod type is GET then methods[j] should be get
					for j := 0; j < ml; j++ {
						if methods[j] == requestMethod {
							methodMatched = true
							break
						}
					}
				}
				// if methodMatched is true then it converts the regexp endpoint to uri form
				// eg: ^/college$ -> /college
				// eg: ^/college/(?P<v0>[^/]+)$ -> /college/{collegeID} (regex is example and not literal conversion)
				if methodMatched {
					// uri will look like -> /college
					uri := r.RequestURI
					// matches the excludedRegex endpoint to the uri
					// eg: (^/college$).MatchString(/college)
					if excludedRoutesRegexp[i].MatchString(strings.Split(uri, "?")[0]) {
						exclude = true
						break
					}
				}
			}
			// executes the middleware if the endpoint is not found in the excludedRoutes slice
			// else redirects it to the specified endpoint which means it was part the excludedRoutes slice
			if !exclude {
				// fmt.Println("=============================================")
				// fmt.Println("Middleware executes ")
				// fmt.Println("=============================================")
				// here ResponseWriter and Request are passed as Middleware will only check
				err := auth.ValidateToken(w, r)
				if err != nil {
					// fmt.Println("SOME ERROR -> ", err)
					web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
					return
				}
				// fmt.Println("No error")
				next.ServeHTTP(w, r)
			} else {
				// fmt.Println("=============================================")
				// fmt.Println("No Middleware Route excluded")
				// fmt.Println("=============================================")
				next.ServeHTTP(w, r)
			}
		})
	}
}

// Middleware will exclude routes on which middleware should not be applied
func (auth *Authentication) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := auth.ValidateToken(w, r)
		if err != nil {
			web.RespondError(w, errors.NewHTTPError(err.Error(), http.StatusBadRequest))
			return
		}
		next.ServeHTTP(w, r)
	})
}
