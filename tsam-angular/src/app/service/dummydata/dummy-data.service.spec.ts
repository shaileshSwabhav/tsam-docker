import { TestBed } from '@angular/core/testing';

import { DummyDataService } from './dummy-data.service';

describe('DummyDataService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: DummyDataService = TestBed.get(DummyDataService);
    expect(service).toBeTruthy();
  });
});
