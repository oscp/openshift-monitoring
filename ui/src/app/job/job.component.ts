import {Component, OnInit} from '@angular/core';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {Subject} from 'rxjs';
import {JobType} from '../shared/job.types';
import {NotificationsService} from 'angular2-notifications';

@Component({
  selector: 'app-job',
  template: ` 
        <h4>Jobs</h4>
        <button class="btn" (click)="newHttpCheck()">Start HttpCheck</button>
        
        <table class="table table-striped">
            <thead class="thead-inverse">
                <tr>
                <th>JobId</th>
                <th>Status</th>
                <th>JobType</th>
                <th>Params</th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                <tr *ngFor="let d of jobs">
                    <td>{{d.JobId}}</td>
                    <td>{{d.JobStatus}}</td>
                    <td>{{d.JobType}}</td>
                    <td>{{d.Params}}</td>
                    <td><button class="btn btn-primary" (click)="stopJob(d.JobId)">Stop</button> </td>
                </tr>
            </tbody>
        </table>
    `
})
export class JobComponent implements OnInit {
  private socket: Subject<any>;
  private jobs: Array<any>;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
    this.socket = socketService.createOrGetWebsocket();
    this.getJobs();
  }

  ngOnInit() {
    this.socket.subscribe(
      message => {
        let data = JSON.parse(message.data);

        console.log(data.WsType, data.Message);

        switch (data.WsType) {
          case SocketType.WS_ALL_JOBS:
            this.jobs = data.Message;
            break;
          case SocketType.WS_NEW_JOB:
            this.notificationService.success("Job created", "JobId: " + data.Message);
            this.getJobs();
            break;
          case SocketType.WS_JOB_STOP:
            this.notificationService.success("Job stopped", "");
            this.getJobs();
        }
      }
    );
  }

  newHttpCheck() {
    this.socket.next({
      WsType: SocketType.WS_NEW_JOB,
      Message: {
        JobType: JobType.JOB_HTTP_CHECK,
        Params: "http://test.ch"
      }
    });
  }

  stopJob(jobId) {
    this.socket.next({
      WsType: SocketType.WS_JOB_STOP,
      Message: {
        JobId: jobId
      }
    })
  }

  getJobs() {
    this.socket.next({WsType: SocketType.WS_ALL_JOBS});
  }

}
