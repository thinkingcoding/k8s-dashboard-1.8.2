
/**
 *
 * @final
 */
export class SessionService {

  constructor() {
    this.isRefreshingSession = false;
  }

  /**
   * Initializes the service.
   */
  init() {
    this.isRefreshingSession = false;
  }

  /**
   *
   * @export
   * @return {boolean}
   */
  isRefreshing(){
    return this.isRefreshingSession;
  }

  /**
   *
   * @export
   */
  startRefresh(){
    this.isRefreshingSession = true;
  }

  /**
   *
   * @export
   */
  endRefresh(){
    this.isRefreshingSession = false;
  }
}
