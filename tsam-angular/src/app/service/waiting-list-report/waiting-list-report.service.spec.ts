import { TestBed } from '@angular/core/testing';

import { WaitingListReportService } from './waiting-list-report.service';

describe('WaitingListReportService', () => {
  let service: WaitingListReportService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WaitingListReportService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
