/** @final */
export class HttpInterceptor {
  /**
   * @param {!angular.$q} $q
   * @param {!angular.$timeout} $timeout
   * @param {!./service.SessionService} kdSessionService
   * @ngInject
   */
  constructor($q, $timeout, kdSessionService) {
    this.request = (config) => {
      // Filter requests made to our backend starting with 'api/v1' and append request header
      // with token stored in a cookie.
      if (config.url.indexOf('api/v1') !== -1) {
        let deferred = $q.defer();
        console.log(config);
        if(kdSessionService.isRefreshing()){
          waitingResolve(deferred, config, [30]);
        }else{
          resolve(deferred, config);
        }
        return deferred.promise;
      }

      return config;
    };
    this.responseError = response => {
      if(response.status === 401){
        location.href = '/logout';
      }
      return $q.reject(response);
    };
    this.response = response => {
      return response;
    };

    function resolve(defer, config){
      defer.resolve(config);
    }
    function waitingResolve(defer, config, $count){
      if($count[0] <= 0) {
        resolve(defer, config);
        return;
      }
      $count[0] = $count[0] - 1;
      $timeout(()=>{
        if(kdSessionService.isRefreshing()){
          waitingResolve(defer, config, $count);
        }else{
          resolve(defer, config);
        }
      }, 1000);
    }
  }

  /**
   *
   * @param {!angular.$q} $q
   * @param {!angular.$timeout} $timeout
   * @param {!./service.SessionService} kdSessionService
   * @return {./interceptor.HttpInterceptor}
   * @ngInject
   */
  static NewHttpInterceptor($q, $timeout, kdSessionService) {
    return new HttpInterceptor($q, $timeout, kdSessionService);
  }
}
