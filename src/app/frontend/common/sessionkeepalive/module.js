import {HttpInterceptor} from './interceptor';
import {SessionService} from './service'
import sessionKeepalvieDirective from './directive';

export default angular
    .module(
        'kubernetesDashboard.sessionkeepalive',
        [])
    .service('kdSessionService', SessionService)
    .factory('kdHttpInterceptor', HttpInterceptor.NewHttpInterceptor)
    .config(initHttpInterceptor)
    .directive('sessionKeepalive', sessionKeepalvieDirective)
    .run(initSessionService);

/**
 *
 * @param {./service.SessionService} kdSessionService
 * @ngInject
 */
function initSessionService(kdSessionService) {
  kdSessionService.init();
}

/**
 * @param {!angular.$HttpProvider} $httpProvider
 * @ngInject
 */
function initHttpInterceptor($httpProvider) {
  $httpProvider.interceptors.push('kdHttpInterceptor');
}
