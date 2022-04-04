import { TestBed } from '@angular/core/testing';

import { TalentDashboardService } from './talent-dashboard.service';

describe('TalentDashboardService', () => {
  let service: TalentDashboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TalentDashboardService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
