import {Component, OnInit} from '@angular/core';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {Subject} from 'rxjs';
import {JobType} from '../shared/job.types';

@Component({
  selector: 'app-job',
  template: ` 
        <h4>Jobs</h4>
        <button (click)="newHttpCheck()">Start HttpCheck</button>
        <button (click)="stopJob()">Stop Job</button>
    `
})
export class JobComponent implements OnInit {
  private socket: Subject<any>;

  constructor(private socketService: SocketService) {
    this.socket = socketService.createOrGetWebsocket();
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

  ngOnInit() {
  }

}
