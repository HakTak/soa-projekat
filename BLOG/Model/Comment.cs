using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;
using System;

namespace BLOG.Model
{
    public class Comment
    {
        [BsonId]
        [BsonRepresentation(BsonType.ObjectId)]
        public string? Id { get; set; } // jedinstveni ID komentara

        [BsonElement("PostId")]
        public string? PostId { get; set; }  // ID blog posta na koji komentar ide

        [BsonElement("AuthorName")]
        public string? AuthorName { get; set; } 

        [BsonElement("AuthorId")]
        public string AuthorId { get; set; } //ID autora komentara

        [BsonElement("Text")]
        public string Text { get; set; }

        [BsonElement("CreatedAt")]
        public DateTime? CreatedAt { get; set; } = DateTime.UtcNow;

        [BsonElement("UpdatedAt")]
        public DateTime? UpdatedAt { get; set; } = DateTime.UtcNow;
    }
}
