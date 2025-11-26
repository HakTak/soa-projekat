import { TestBed } from '@angular/core/testing';

import { ProfileFollowService } from './profile-follow.service';

describe('ProfileFollowService', () => {
  let service: ProfileFollowService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProfileFollowService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
