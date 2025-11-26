import { TestBed } from '@angular/core/testing';

import { SimulatorServiceService } from './simulator-service.service';

describe('SimulatorServiceService', () => {
  let service: SimulatorServiceService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(SimulatorServiceService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
