import { TestBed } from '@angular/core/testing';

import { CareerObjectiveService } from './career-objective.service';

describe('CareerObjectiveService', () => {
  let service: CareerObjectiveService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(CareerObjectiveService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
