import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {FormsModule} from '@angular/forms';
import {AppComponent} from './app.component';
import {SocketService} from './socket.service';
import {SimpleNotificationsModule} from 'angular2-notifications';
import {DaemonsComponent} from './daemons/daemons.component';
import { ChecksComponent } from './checks/checks.component';
import { ResultsComponent } from './results/results.component';
import {ChartsModule} from "ng2-charts";
import {NotificationsService} from "angular2-notifications";
import {HttpClientModule} from "@angular/common/http";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

@NgModule({
  declarations: [
    AppComponent,
    DaemonsComponent,
    ChecksComponent,
    ResultsComponent
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    FormsModule,
    HttpClientModule,
    SimpleNotificationsModule,
    ChartsModule
  ],
  providers: [SocketService, NotificationsService],
  bootstrap: [AppComponent]
})
export class AppModule {
}
