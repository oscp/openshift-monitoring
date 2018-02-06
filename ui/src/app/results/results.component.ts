import {Component, OnInit, ViewChild} from '@angular/core';
import {NotificationsService} from "angular2-notifications";
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";
import {BaseChartDirective} from "ng2-charts";

@Component({
  selector: 'app-results',
  templateUrl: 'results.component.html'
})
export class ResultsComponent implements OnInit {
  // Dognut charts
  public dognutChartType: string = 'doughnut';
  public dognutChartOptions: any = {
    legend: {
      display: false
    }
  };
  public checkTypeLabels: string[] = ["MASTER_API_CHECK", "DNS_NSLOOKUP_KUBERNETES", "DNS_SERVICE_NODE",
    "DNS_SERVICE_POD", "HTTP_POD_SERVICE_A_B", "HTTP_POD_SERVICE_A_C", "HTTP_SERVICE_ABC", "HTTP_HAPROXY", "ETCD_HEALTH"];

  public errorData: number[] = [0, 0, 0, 0, 0, 0, 0, 0, 0];
  public successData: number[] = [0, 0, 0, 0, 0, 0, 0, 0, 0];
  public errors: Array<any> = [];

  public checkOverviewLabels: string[] = ['Started', 'Finished'];
  public checkOverviewData: number[] = [0, 0];

  // Line Chart
  @ViewChild('linechart') chart: BaseChartDirective;
  public lineChartType: string = 'line';
  public checkLineData: any = [
    {data: [], label: 'Successful checks'},
    {data: [], label: 'Failed checks'}
  ];
  public checkLineLabels: Array<any> = [];
  public checkLineLegend: boolean = true;
  public checkLineOptions: any = {
    responsive: true
  };
  private successCount: number = 0;
  private errorCount: number = 0;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);
        switch (data.type) {
          case SocketType.CHECK_RESULTS:
            if (data.message.length > 0) {
              this.handleResults(data.message);
            }
            break;
          case SocketType.ALL_DAEMONS:
            this.handleDaemonUpdate(data.message);
            break;
        }
      }
    );
  }

  private handleDaemonUpdate(daemons) {
    this.checkOverviewData[0] = 0;
    this.checkOverviewData[1] = 0;

    daemons.forEach(d => {
      this.checkOverviewData[0] += d.startedChecks;
      this.checkOverviewData[1] += d.failedChecks + d.successfulChecks;
    });

    // Force UI update
    this.checkOverviewData = this.checkOverviewData.slice();
  }

  private handleResults(msg) {
    msg.forEach(m => {
      // Handle specific by result
      if (m.isOk) {
        this.handleSuccessResult(m);
      } else {
        this.handleErrorResult(m);
      }
    });

    // Handle Line-Charts
    this.handleLineResult();
  }

  private handleLineResult() {
    let now = new Date();
    this.checkLineLabels.push(`${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`);
    this.checkLineData[0].data.push(this.successCount);
    this.checkLineData[1].data.push(this.errorCount);

    // Cleanup counters
    this.successCount = 0;
    this.errorCount = 0;

    // Update UI because of bug in chartjs:
    this.chart.labels = this.checkLineLabels.slice();
    this.checkLineData = this.checkLineData.slice();
  }

  private handleErrorResult(msg) {
    this.errorCount++;
    let idx = this.checkTypeLabels.findIndex(m => m == msg.type);

    if (idx > -1) {
      this.errorData[idx] += 1;
    }

    // Enforce refresh
    this.errorData = this.errorData.slice();

    // Tell the user about it
    msg.date = new Date();
    this.errors.unshift(msg);
  }

  private handleSuccessResult(msg) {
    this.successCount++;
    let idx = this.checkTypeLabels.findIndex(m => m == msg.type);

    if (idx > -1) {
      this.successData[idx] += 1;
    }

    // Enforce refresh
    this.successData = this.successData.slice();
  }
}
