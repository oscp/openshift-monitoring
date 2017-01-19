import { OseMonitoringUiPage } from './app.po';

describe('ose-monitoring-ui App', function() {
  let page: OseMonitoringUiPage;

  beforeEach(() => {
    page = new OseMonitoringUiPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
