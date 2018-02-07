import {Component, OnInit, SimpleChanges, ViewChild} from '@angular/core';
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

  public checkOverviewLabels: string[] = ['Started', 'Finished'];
  public checkOverviewData: number[] = [0, 0];

  public checkTypeLabels: string[] = ["MASTER_API_CHECK", "DNS_NSLOOKUP_KUBERNETES", "DNS_SERVICE_NODE",
    "DNS_SERVICE_POD", "HTTP_POD_SERVICE_A_B", "HTTP_POD_SERVICE_A_C", "HTTP_SERVICE_ABC", "HTTP_HAPROXY", "ETCD_HEALTH"];
  public errorData: number[] = [0, 0, 0, 0, 0, 0, 0, 0, 0];
  public successData: number[] = [0, 0, 0, 0, 0, 0, 0, 0, 0];

  public failures: Array<any> = [];

  // Line Chart
  @ViewChild(BaseChartDirective) chart: BaseChartDirective;
  public lineChartType: string = 'line';
  public checkLineData: any = [
    {data: [], label: 'Successful checks'},
    {data: [], label: 'Failed checks'}
  ];
  public checkLineLabels: Array<any> = [];
  public checkLineLegend = true;
  public checkLineOptions: any = {
    responsive: true
  };

  constructor(private socketService: SocketService) {
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);
        switch (data.type) {
          case SocketType.CHECK_RESULTS:
            this.handleResults(data.message);
            break;
        }
      }
    );
  }

  private handleResults(res) {
    // Failures
    this.failures = res.failures.slice().reverse();

    // Started & finished checks
    this.checkOverviewData[0] = res.startedChecks;
    this.checkOverviewData[1] = res.finishedChecks;
    this.checkOverviewData = this.checkOverviewData.slice();

    // Success / failed by type
    this.handleFailedByType(res.failedChecksByType);
    this.handleSuccesfulByType(res.successfulChecksByType);

    // Handle Line-Charts
    this.handleLineResult(res);
  }

  private handleLineResult(res: any) {
    this.checkLineLabels = [];
    this.checkLineData[0].data = [];
    this.checkLineData[1].data = [];

    for (let [k, v] of Object.entries(res.ticks)) {
      this.checkLineLabels.push(k);
      this.checkLineData[0].data.push(v.successfulChecks);
      this.checkLineData[1].data.push(v.failedChecks);
    }
    this.chart.chart.update();
  }

  private handleFailedByType(res: any) {
    for (let [k, v] of Object.entries(res)) {
      // Find index for key
      let idx = this.checkTypeLabels.findIndex(m => m === k);
      this.errorData[idx] = v;
    }

    // Enforce refresh
    this.errorData = this.errorData.slice();
  }

  private handleSuccesfulByType(res: any) {
    for (let [k, v] of Object.entries(res)) {
      // Find index for key
      let idx = this.checkTypeLabels.findIndex(m => m === k);
      this.successData[idx] = v;
    }

    // Enforce refresh
    this.successData = this.successData.slice();
  }
}
