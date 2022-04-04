import { TestBed } from '@angular/core/testing';

import { ConceptDashboardService } from './concept-dashboard.service';

describe('ConceptDashboardService', () => {
  let service: ConceptDashboardService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ConceptDashboardService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
