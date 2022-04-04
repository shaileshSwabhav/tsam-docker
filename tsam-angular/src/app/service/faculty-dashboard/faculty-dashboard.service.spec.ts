import { TestBed } from '@angular/core/testing';

import { FacultyDashboardService } from './faculty-dashboard.service';

describe('FacultyDashboardService', () => {
  let service: FacultyDashboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FacultyDashboardService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
