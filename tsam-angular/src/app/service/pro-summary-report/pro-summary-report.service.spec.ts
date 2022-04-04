import { TestBed } from '@angular/core/testing';

import { ProSummaryReportService } from './pro-summary-report.service';

describe('ProSummaryReportService', () => {
  let service: ProSummaryReportService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProSummaryReportService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
