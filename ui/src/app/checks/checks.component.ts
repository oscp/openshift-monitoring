import {Component, OnInit} from '@angular/core';
import {SocketService} from "../socket.service";
import {SocketType} from "../shared/socket.types";

@Component({
  selector: 'app-checks',
  templateUrl: 'checks.component.html'
})
export class ChecksComponent implements OnInit {
  public checks = {};

  constructor(private socketService: SocketService) {
    this.getCurrentChecks();
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);
        switch (data.type) {
          case SocketType.CURRENT_CHECKS:
            this.checks = data.message;
            break;
        }
      }
    );
  }

  public startChecks() {
    this.socketService.websocket.next({type: SocketType.START_CHECKS, message: this.checks});
  }

  public stopChecks() {
    this.socketService.websocket.next({type: SocketType.STOP_CHECKS});
  }

  public resetStats() {
    this.socketService.websocket.next({type: SocketType.RESET_STATS});
  }

  private getCurrentChecks() {
    this.socketService.websocket.next({Type: SocketType.CURRENT_CHECKS});
  }
}
