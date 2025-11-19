using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace BLOG.Model
{
    public class PostLike
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string Id {get; set;} = null!;

        [BsonElement("UserId")]
        public string UserId { get; set; } = null!;
 
        [BsonElement("PostId")]
        public string? PostId { get; set; } = null!;

        [BsonElement("LikedAt")]
        public DateTime LikedAt { get; set; } = DateTime.UtcNow;

    }
}