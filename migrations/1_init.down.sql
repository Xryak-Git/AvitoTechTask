drop table if exists organization_responsible cascade;

drop table if exists tender_versions cascade;

drop table if exists bid_bidreview cascade;

drop table if exists bid_review cascade;

drop table if exists bid_versions cascade;

drop table if exists decision cascade;

drop table if exists bid cascade;

drop table if exists employee cascade;

drop table if exists tender cascade;

drop table if exists organization cascade;

drop function if exists save_tender_version() cascade;

drop function if exists save_and_increment_tender_version() cascade;

drop function if exists save_and_increment_bid_version() cascade;

drop function if exists update_bid_status() cascade;

drop type if exists organization_type cascade;

drop type if exists service_type cascade;

drop type if exists tender_status cascade;

drop type if exists bid_status cascade;

drop type if exists author_type cascade;

