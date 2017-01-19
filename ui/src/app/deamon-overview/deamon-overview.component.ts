import { Component, OnInit } from '@angular/core';
import {Subject} from "rxjs";
import {SocketService} from "../socket.service";

@Component({
  selector: 'app-deamon-overview',
  template: `
        <p>hi</p>
    `
})
export class DeamonOverviewComponent implements OnInit {
  private socket: Subject<any>;

  constructor(private socketService: SocketService) {
    this.socket = socketService.createOrGetWebsocket();
  }

  ngOnInit() {
    this.socket.subscribe(
        message => {
          let data = JSON.parse(message.data);
          console.log(data);
        }
    );
  }
}
