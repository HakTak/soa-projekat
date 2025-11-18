
using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace BLOG.Model
{
    public class Post
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id {get; set;}

        [BsonElement("Title")]
        public string Title {get; set;}

        [BsonElement("Description")]
        public string Description {get; set;}

        [BsonElement("CreatedAt")]
        public DateTime CreatedAt { get; set; }

        [BsonElement("ImagePaths")]
        public string[] ImagePaths {get; set;}

        [BsonElement("LikeCount")]
        public int LikeCount {get; set;}

    }
}