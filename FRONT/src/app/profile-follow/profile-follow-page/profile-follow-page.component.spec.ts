import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProfileFollowPageComponent } from './profile-follow-page.component';

describe('ProfileFollowPageComponent', () => {
  let component: ProfileFollowPageComponent;
  let fixture: ComponentFixture<ProfileFollowPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProfileFollowPageComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ProfileFollowPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
